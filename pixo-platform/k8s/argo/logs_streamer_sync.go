package argo

import "github.com/rs/zerolog/log"

func (s *LogsStreamer) IsDone() bool {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	isDone := s.numNodes == s.numDone

	log.Debug().Msgf("is completely done: %t", isDone)
	return isDone
}

func (s *LogsStreamer) NumNodes() int {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	log.Debug().Msgf("num total nodes: %d", s.numNodes)
	return s.numNodes
}

func (s *LogsStreamer) NumDone() int {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	log.Debug().Msgf("num nodes done: %d", s.numDone)
	return s.numDone
}

func (s *LogsStreamer) addStream(name string) chan Log {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if s.streams[name] == nil {
		log.Debug().Msgf("opening new stream for node %s", name)
		s.streams[name] = make(chan Log, 100)
		s.numNodes++
	}

	return s.streams[name]
}

func (s *LogsStreamer) getStream(name string) chan Log {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	log.Debug().Msgf("getting stream for node %s", name)
	return s.streams[name]
}

func (s *LogsStreamer) markStreamDone(name string) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.numDone++
	close(s.streams[name])

	log.Debug().Msgf("marked stream done for node %s", name)
}
