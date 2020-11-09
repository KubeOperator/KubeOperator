import {Component, OnInit} from '@angular/core';
import {LicenseService} from '../../business/setting/license/license.service';
import {SessionService} from "../../shared/auth/session.service";
import {SessionUser} from "../../shared/auth/session-user";

@Component({
    selector: 'app-navigation',
    templateUrl: './navigation.component.html',
    styleUrls: ['./navigation.component.css']
})
export class NavigationComponent implements OnInit {

    constructor(private licenseService: LicenseService, private sessionService: SessionService) {
    }

    hasLicense = false;
    user: SessionUser;

    ngOnInit(): void {
        this.licenseService.get().subscribe(data => {
            if (data.status === 'valid') {
                this.hasLicense = true;
            }
        });
        const profile = this.sessionService.getCacheProfile();
        this.user = profile.user;
    }

}
