package helpers

import (
	"github.com/AskatNa/OnlineClinic/api/models"
)

func GenerateAppointmentSlots(slotsNo int, doctorId string) []models.AppointmentSlot {
	appointmentSlots := make([]models.AppointmentSlot, slotsNo)

	for i := range appointmentSlots {
		appointmentSlots[i].SlotNo = i + 1
	}

	return appointmentSlots
}

func GenerateWeekDoctorSchedule(doctorId string) models.DoctorSchedule {

	weekDays := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	schedule := make(map[string]models.ScheduleDay)

	appointmentSlots := GenerateAppointmentSlots(3, doctorId)

	for _, day := range weekDays {
		sd := models.ScheduleDay{
			AppointmentSlots: appointmentSlots,
		}
		schedule[day] = sd
	}

	return models.DoctorSchedule{
		WeeklySchedule: schedule,
	}
}
