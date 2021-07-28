package server

import "sync"

type ClientManager struct {
	sync.Mutex
	clients map[string]*Client
	users   map[string][]string
}

func NewClientManager() *ClientManager {
	cm := new(ClientManager)
	cm.clients = make(map[string]*Client)
	cm.users = make(map[string][]string)
	return cm
}

func (cm *ClientManager) Add(client *Client) {
	cm.Lock()
	cm.clients[client.ConnId] = client
	users := cm.users[client.Uin]
	users = append(users, client.ConnId)
	cm.Unlock()
}

func (cm *ClientManager) Get(id string) *Client {
	cm.Lock()
	defer cm.Unlock()

	client := cm.clients[id]
	return client
}

func (cm *ClientManager) Remove(id string) (client *Client) {
	cm.Lock()
	defer cm.Unlock()

	client = cm.clients[id]
	if client != nil {
		delete(cm.clients, id)

		users := cm.users[client.Uin]
		for i, v := range users {
			if v == client.ConnId {
				users = append(users[:i], users[i+1:]...)
				cm.users[client.Uin] = users
				break
			}
		}

	}

	return
}

func (cm *ClientManager) GetUserClients(uin string) []*Client {
	cm.Lock()
	defer cm.Unlock()

	users := cm.users[uin]
	clients := make([]*Client, 0, len(users))
	for _, id := range users {
		client := cm.clients[id]
		if client != nil {
			clients = append(clients, client)
		}
	}

	return clients
}

func (cm *ClientManager) GetClient(uin, plat string) *Client {
	clients := cm.GetUserClients(uin)

	for _, client := range clients {
		if client.Platform == plat {
			return client
		}
	}

	return nil
}
