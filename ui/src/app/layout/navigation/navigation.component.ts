import {Component, OnInit} from '@angular/core';
import {SessionService} from "../../shared/auth/session.service";
import {SessionUser} from "../../shared/auth/session-user";
import {BusinessLicenseService} from '../../shared/service/business-license.service';

@Component({
    selector: 'app-navigation',
    templateUrl: './navigation.component.html',
    styleUrls: ['./navigation.component.css']
})
export class NavigationComponent implements OnInit {

    constructor(private businessLicenseService: BusinessLicenseService, private sessionService: SessionService) {
    }

    hasLicense = false;
    user: SessionUser;

    ngOnInit(): void {
        this.hasLicense = this.businessLicenseService.licenseValid;
        this.sessionService.getProfile().subscribe(res => {
            this.user = res.user;
        })
    }

}
