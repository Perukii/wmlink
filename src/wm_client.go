package main

func (host *WmHost) wm_client_allocate_from_host() WmClientAddress{
	var clt WmClient
	host.client = append(host.client, clt)
	return len(host.client)-1
}

func (host *WmHost) wm_client_search(window XWindowID) WmClientAddress{
	for i := 0; i < len(host.client); i++ {
		if host.client[i].box.window == window { return i }
		if host.client[i].app == window { return i }
		if host.client[i].mask.window == window { return i }
	}

	return 0
}

func (host *WmHost) wm_client_withdraw(address WmClientAddress, app_is_destroyed bool){

	clt := host.client[address]
	attr := host.wm_host_get_window_attributes(clt.box.window)
	attr.override_redirect = attr.override_redirect
	
	if app_is_destroyed == false{
		host.wm_host_reparent_window(clt.app, host.wm_host_query_parent(clt.box.window),
			int(attr.x), int(attr.y))
	}

	host.wm_host_remove_transparent(host.client[address].box)
	host.wm_host_remove_transparent(host.client[address].mask)

	host.client[address] = host.client[len(host.client)-1]
	host.client = host.client[:len(host.client)-1]
	
}

func (host *WmHost) wm_client_setup(clt *WmClient, xapp XWindowID){

	clt.app = xapp

	attr := host.wm_host_get_window_attributes(clt.app)

	border_width := host.config.client_drawable_range_border_width

	// ---Box---
	host.wm_host_setup_transparent(&clt.box, host.root_window, int(attr.x), int(attr.y),
								   int(attr.width) +border_width*2,
								   int(attr.height)+border_width*2)
	clt.box.drawtype = WM_DRAW_TYPE_BOX

	host.wm_host_select_input(clt.box.window, XSubstructureNotifyMask | XSubstructureRedirectMask)
	host.wm_host_map_window(clt.box.window)
	host.wm_host_raise_window(clt.box.window)
	host.wm_host_reparent_window(clt.app, clt.box.window, border_width, border_width)

	host.wm_host_draw_transparent(clt.box)


	box_attr := host.wm_host_get_window_attributes(clt.box.window)

	// ---Mask---
	clt.mask.drawtype = WM_DRAW_TYPE_MASK
	host.wm_host_setup_transparent(&clt.mask, host.root_window,
								   int(box_attr.x),
								   int(box_attr.y),
								   int(box_attr.width),
								   int(box_attr.height))

	host.wm_host_select_input(clt.mask.window, XSubstructureNotifyMask)
	host.wm_host_map_window(clt.mask.window)
}

func (host *WmHost) wm_client_deactivate(address WmClientAddress){
	clt := host.client[address]
	attr := host.wm_host_get_window_attributes(clt.box.window)
	host.wm_client_configure(address, int(attr.x), int(attr.y), int(attr.width), int(attr.height))
	host.wm_host_raise_window(clt.box.window)
	host.wm_host_raise_window(clt.mask.window)
	host.wm_host_draw_transparent(clt.mask)
}

func (host *WmHost) wm_client_activate(address WmClientAddress){
	clt := host.client[address]
	host.wm_host_lower_window(clt.mask.window)
	host.wm_host_raise_window(clt.box.window)
}

func (host *WmHost) wm_client_configure(address WmClientAddress, x int, y int, w int, h int){
	clt := host.client[address]
	host.wm_host_move_window  (clt.box.window, x, y)
	host.wm_host_resize_window(clt.box.window, w, h)
	host.wm_host_draw_transparent(clt.box)

	host.wm_host_move_window  (clt.mask.window, x, y)
	host.wm_host_resize_window(clt.mask.window, w, h)
	host.wm_host_draw_transparent(clt.mask)
}

func (host *WmHost) wm_client_adjust_mask_for_box(address WmClientAddress){
	clt := host.client[address]
	attr := host.wm_host_get_window_attributes(clt.box.window)
	host.wm_host_move_window  (clt.mask.window, int(attr.x), int(attr.y))
	host.wm_host_resize_transparent(clt.mask, int(attr.width), int(attr.height))
	host.wm_host_draw_transparent(clt.mask)
}

func (host *WmHost) wm_client_adjust_app_for_box(address WmClientAddress){
	clt := host.client[address]
	attr := host.wm_host_get_window_attributes(clt.box.window)
	border_width := host.config.client_drawable_range_border_width
	host.wm_host_resize_window(clt.app, int(attr.width)-border_width*2, int(attr.height)-border_width*2)
}