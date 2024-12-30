package router

func commonGroups() []CommonRouter {
	return []CommonRouter{
		&BaseRouter{},
		&DashboardRouter{},
		&HostRouter{},
		&ContainerRouter{},
		&LogRouter{},
		&FileRouter{},
		&ToolboxRouter{},
		&TerminalRouter{},
		&SettingRouter{},
		&WebsiteGroupRouter{},
		&ProcessRouter{},
	}
}
