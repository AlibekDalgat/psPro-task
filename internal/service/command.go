package service

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"os/exec"
	"psPro-task/internal/models"
	"psPro-task/internal/repository"
	"syscall"
	"time"
)

const writeInterval = 3

type CommandService struct {
	repo        repository.Command
	currentJobs *map[int]models.Job
}

func NewCommandService(repo *repository.Repository) *CommandService {
	return &CommandService{repo: repo, currentJobs: &map[int]models.Job{}}
}

func (s *CommandService) CreateCommand(command models.Command) (int, error) {
	command.CreatedAt = time.Now()
	command.ExecutedAT = nil
	command.Stdout = nil
	command.Stderr = nil

	id, err := s.repo.CreateCommand(command)
	if err != nil {
		return 0, err
	}
	(*s.currentJobs)[id] = models.Job{ActionChan: make(chan string)}
	return id, nil
}

func (s *CommandService) ExecuteCommand(id int, script string) {
	cmd := exec.Command("bash", "-c", script)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logrus.Fatal(err)
		return
	}
	stdoutScanner := bufio.NewScanner(stdout)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		logrus.Fatal(err)
		return
	}
	stderrScanner := bufio.NewScanner(stderr)

	if err = cmd.Start(); err != nil {
		logrus.Printf("Не получилось запустить команду %d: '%s'\n", id, err.Error())
		delete(*s.currentJobs, id)
	} else {
		job := (*s.currentJobs)[id]
		job.IsRun = true
		job.StdOutBuffer = new(bytes.Buffer)
		job.StdErrBuffer = new(bytes.Buffer)
		(*s.currentJobs)[id] = job
	}

	go s.scanStdStream(stdoutScanner, (*s.currentJobs)[id].StdOutBuffer)
	go s.scanStdStream(stderrScanner, (*s.currentJobs)[id].StdErrBuffer)

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	for {
		select {
		case <-done:
			logrus.Printf("Команда %d завершена\n", id)
			outbuf := (*s.currentJobs)[id].StdOutBuffer.String()
			errbuf := (*s.currentJobs)[id].StdErrBuffer.String()
			err = s.repo.WriteResults(id, &outbuf, &errbuf, time.Now())
			if err != nil {
				logrus.Printf("Запись результатов выполнения команды %d не получилось: '%s'\n", id, err.Error())
			}
			delete(*s.currentJobs, id)
			break
		case action := <-(*s.currentJobs)[id].ActionChan:
			switch action {
			case "stop":
				err = cmd.Process.Signal(syscall.SIGSTOP)
				if err != nil {
					logrus.Printf("Не получилось отправить сигнал остановки команде с id %d: '%s'\n", id, err.Error())
				}
				logrus.Printf("Команда %d остановлена", id)
				job := (*s.currentJobs)[id]
				job.IsRun = false
				(*s.currentJobs)[id] = job
			case "start":
				err = cmd.Process.Signal(syscall.SIGCONT)
				if err != nil {
					logrus.Printf("Не получилось отправить сигнал старта команде с id %d: '%s'\n", id, err.Error())
				}
				logrus.Printf("Команда %d продолжена", id)
				job := (*s.currentJobs)[id]
				job.IsRun = true
				(*s.currentJobs)[id] = job
			case "kill":
				if err = cmd.Process.Kill(); err != nil {
					logrus.Errorf("ошибка завершения команды %d: %s", id, err.Error())
				}
				logrus.Printf("Команда %d принудительна завершена", id)
				break
			}
		}
	}
}

func (s *CommandService) scanStdStream(scanner *bufio.Scanner, buffer *bytes.Buffer) {
	for scanner.Scan() {
		buffer.WriteString(scanner.Text())
		buffer.WriteString("\n")
	}
}

func (s *CommandService) StopCommand(id int) error {
	job, ok := (*s.currentJobs)[id]
	if !ok {
		return fmt.Errorf("Команды с id %d не найдено", id)
	}
	if !job.IsRun {
		return fmt.Errorf("Команда %d уже остановлена", id)
	}
	job.ActionChan <- "stop"
	return nil
}

func (s *CommandService) StartCommand(id int) error {
	job, ok := (*s.currentJobs)[id]
	if !ok {
		return fmt.Errorf("Команды с id %v не найдено", id)
	}
	if job.IsRun {
		return fmt.Errorf("Команда %d уже запущена", id)
	}
	job.ActionChan <- "start"
	return nil
}

func (s *CommandService) KillCommand(id int) error {
	job, ok := (*s.currentJobs)[id]
	if !ok {
		return fmt.Errorf("Команды с id %v не найдено", id)
	}
	job.ActionChan <- "kill"
	return nil
}

func (s *CommandService) GetAllCommands() ([]models.Command, error) {
	return s.repo.GetAllCommands()
}

func (s *CommandService) GetOneCommand(id int) (models.Command, error) {
	command, err := s.repo.GetOneCommand(id)
	if err != nil {
		return models.Command{}, err
	}
	if job, ok := (*s.currentJobs)[id]; ok {
		command.Stdout = new(string)
		*command.Stdout = job.StdOutBuffer.String()
		command.Stderr = new(string)
		*command.Stderr = job.StdErrBuffer.String()
	}
	return command, nil
}
