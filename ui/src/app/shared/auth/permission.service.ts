import {Injectable} from '@angular/core';
import {SessionService} from './session.service';
import {Profile} from './session-user';

@Injectable({
    providedIn: 'root'
})
export class PermissionService {


    constructor(private sessionService: SessionService) {
    }

    async authRootMenu(menu: string) {
        let userProfile = this.sessionService.getCacheProfile();
        if (userProfile == null) {
            userProfile = await this.getProfile();
        }
        if (userProfile.user.isAdmin) {
            return true;
        }
        for (const roleMenu of userProfile.roleMenus) {
            for (const me of roleMenu.menus) {
                if (me === menu) {
                    return true;
                }
            }
        }
        return false;
    }

    async authSecondaryMenu(menu: string, projectId: string) {
        let userProfile = this.sessionService.getCacheProfile();
        if (userProfile == null) {
            userProfile = await this.getProfile();
        }
        if (userProfile.user.isAdmin) {
            return true;
        }
        for (const roleMenu of userProfile.roleMenus) {
            if (roleMenu.projectId === projectId) {
                for (const me of roleMenu.menus) {
                    if (me === menu) {
                        return true;
                    }
                }
            }

        }
        return false;
    }

    async authOperate(operate: string, projectName: string) {
        const auths = operate.split('.');
        if (auths.length < 2) {
            return false;
        }
        let userProfile = this.sessionService.getCacheProfile();
        if (userProfile == null) {
            userProfile = await this.getProfile();
        }
        if (userProfile.user.isAdmin) {
            return true;
        }
        let result = false;
        const op = auths[0];
        const ro = auths[1];

        start:
            for (const permission of userProfile.permissions) {
                if (permission.projectName === projectName) {
                    for (const userPermissionRole of permission.userPermissionRoles) {
                        if (userPermissionRole.operation === op) {
                            for (const role of userPermissionRole.roles) {
                                if (role === ro) {
                                    result = true;
                                    break start;
                                }
                            }
                        }
                    }
                }
            }
        return result;
    }

    authOp(operate: string, projectName: string): boolean {
        let result = false;
        const auths = operate.split('.');
        if (auths.length < 2) {
            return false;
        }
        const userProfile = this.sessionService.getCacheProfile();
        if (userProfile == null) {
            return false;
        }
        if (userProfile.user.isAdmin) {
            return true;
        }
        const op = auths[0];
        const ro = auths[1];
        start:
            for (const permission of userProfile.permissions) {
                if (permission.projectName === projectName) {
                    for (const userPermissionRole of permission.userPermissionRoles) {
                        if (userPermissionRole.operation === op) {
                            for (const role of userPermissionRole.roles) {
                                if (role === ro) {
                                    result = true;
                                    break start;
                                }
                            }
                        }
                    }
                }
            }
        return result;
    }

    getProjectRole(projectName: string) {
        const userProfile = this.sessionService.getCacheProfile();
        if (userProfile == null) {
            return '';
        }
        if (userProfile.user.isAdmin) {
            return 'SYSTEM_ADMIN';
        }
        let projectRole = '';
        for (const permission of userProfile.permissions) {
            if (permission.projectName === projectName) {
                projectRole = permission.projectRole;
                break;
            }
        }
        return projectRole;
    }

    getProfile() {
        return this.sessionService.getProfile().toPromise();
    }
}
