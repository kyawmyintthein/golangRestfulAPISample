package mongo

// import "gopkg.in/mgo.v2"

// type Service struct {
// 	baseSession *mgo.Session
// 	queue       chan int
// 	URL         string
// 	Open        int
// }

// var service Service

// func (s *Service) New() error {
// 	var err error
// 	s.queue = make(chan int, maxPool)
// 	for i := 0; i < maxPool; i = i + 1 {
// 		s.queue <- 1
// 	}
// 	s.Open = 0
// 	s.baseSession, err = mgo.Dial(s.URL)
// 	return err
// }

// func (s *Service) Session() *mgo.Session {
// 	<-s.queue
// 	s.Open++
// 	return s.baseSession.Copy()
// }

// func (s *Service) Close(c *Collection) {
// 	c.db.s.Close()
// 	s.queue <- 1
// 	s.Open--
// }
