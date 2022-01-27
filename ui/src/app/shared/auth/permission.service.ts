import {Injectable} from '@angular/core';
import {SessionService} from './session.service';

@Injectable({
    providedIn: 'root'
})
export class PermissionService {


    constructor(private sessionService: SessionService) {
    }

    async authRootMenu(menu: string) {
        this.sessionService.getProfile().subscribe(res => {
            if (res.user.isAdmin) {
                return true;
            }
            for (const roleMenu of res.roleMenus) {
                for (const me of roleMenu.menus) {
                    if (me === menu) {
                        return true;
                    }
                }
            }
        })
        return false;
    }

    async authSecondaryMenu(menu: string, projectId: string) {
        this.sessionService.getProfile().subscribe(res => {
            if (res.user.isAdmin) {
                return true;
            }
            for (const roleMenu of res.roleMenus) {
                if (roleMenu.projectId === projectId) {
                    for (const me of roleMenu.menus) {
                        if (me === menu) {
                            return true;
                        }
                    }
                }
            }
        })
        return false;
    }

    async authOperate(operate: string, projectName: string) {
        const auths = operate.split('.');
        if (auths.length < 2) {
            return false;
        }
        let profile
        this.sessionService.getProfile().subscribe(res => {
            profile = res
            if (res.user.isAdmin) {
                return true;
            }
        }, error => {
            return false
        })
        let result = false;
        const op = auths[0];
        const ro = auths[1];

        start:
            for (const permission of profile.permissions) {
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
        let profile
        this.sessionService.getProfile().subscribe(res => {
            profile = res
            if (res.user.isAdmin) {
                return true;
            }
        }, error => {
            return false
        })
        const op = auths[0];
        const ro = auths[1];
        start:
            for (const permission of profile.permissions) {
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
        let profile
        this.sessionService.getProfile().subscribe(res => {
            profile = res
            if (res.user.isAdmin) {
                return 'SYSTEM_ADMIN';
            }
        }, error => {
            return ''
        })
        let projectRole = '';
        for (const permission of profile.permissions) {
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
