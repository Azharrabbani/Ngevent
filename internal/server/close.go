package server

import "log"

func (s *FiberServer) Close() {
	if s.App != nil {
		_ = s.App.Shutdown()
	}

	// Close database
	if s.DB != nil {
		sqlDB, err := s.DB.DB()
		if err != nil {
			log.Println("failed to get sql.DB:", err)
		} else if err := sqlDB.Close(); err != nil {
			log.Println("Failed to close database:", err)
		} else {
			log.Println("Database connection closed")
		}
	}

	// Close redis
	if s.RDB != nil {
		if err := s.RDB.Close(); err != nil {
			log.Println("Failed to close redis:", err)
		} else {
			log.Println("Redis connection closed")
		}
	}

	// Close asynq
	if s.ClientWoker != nil {
		s.ClientWoker.Close()
		log.Println("Asynq client closed")
	}

	if s.InspectorWorker != nil {
		s.InspectorWorker.Close()
		log.Println("Asynq inspector closed")
	}
}
