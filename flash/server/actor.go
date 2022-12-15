package server

type userActor struct {
	session Session
}

type PlayerActor struct {
	userActor
}
