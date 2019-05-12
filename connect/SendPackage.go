package connect

import (
    "../models"
    "../settings"
)

func SendPackage(pack models.PackageTCP, mode settings.ModeNet) {
	switch mode {
		case settings.P2P_mode: 
			createRedirectP2PPackage(&pack)
        	sendInitRedirectP2PPackage(pack)
		case settings.F2F_mode:
			sendEncryptedPackage(pack, settings.F2F_mode)
	}
}
