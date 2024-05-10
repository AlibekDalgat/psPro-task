package service

import (
	"bufio"
	"fmt"
	"github.com/sirupsen/logrus"
	"os/exec"
	"psPro-task/internal/models"
	"psPro-task/internal/repository"
	"syscall"
	"time"
)

type CommandService struct {
	repo        repository.Command
	currentJobs *map[int]models.Job
}

func NewCommandService(repo *repository.Repository) *CommandService {
	return &CommandService{repo: repo, currentJobs: &map[int]models.Job{}}
}

func (s *CommandService) CreateCommand(command models.Command) (int, error) {
	command.CreatedAt = time.Now()
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
		(*s.currentJobs)[id] = job
	}

	go func() {
		for stdoutScanner.Scan() {
			line := stdoutScanner.Text()
			err = s.repo.WriteToColumn("stdout", id, line)
			if err != nil {
				logrus.Printf("Запись в stdout команды %d не получилась: '%s'\n", id, err.Error())
			}
		}
	}()

	go func() {
		for stderrScanner.Scan() {
			line := stderrScanner.Text()
			err = s.repo.WriteToColumn("stdout", id, line)
			if err != nil {
				logrus.Printf("Запись в stderr команды %d не получилась: '%s'\n", id, err.Error())
			}
		}
	}()
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-done:
		logrus.Println("process finished successfully")
		err = s.repo.WriteToColumn("executed_at", id, time.Now())
		if err != nil {
			logrus.Printf("Запись в executed_at не получилась: '%s'\n", err.Error())
		}
		delete(*s.currentJobs, id)
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
			delete(*s.currentJobs, id)
			logrus.Printf("Команда %d принудительна завершена", id)
		}
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
	err := s.repo.WriteToColumn("executed_at", id, time.Now())
	if err != nil {
		logrus.Printf("Запись в executed_at команды %d не получилась: '%s'\n", id, err.Error())
	}
	return nil
}
