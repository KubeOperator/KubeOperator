import {Component, OnInit} from '@angular/core';
import {BusinessLicenseService} from '../../shared/service/business-license.service';

@Component({
    selector: 'app-setting',
    templateUrl: './setting.component.html',
    styleUrls: ['./setting.component.css']
})
export class SettingComponent implements OnInit {

    constructor(private businessLicenseService: BusinessLicenseService) {
    }

    hasLicense = false;

    ngOnInit(): void {
        this.hasLicense = this.businessLicenseService.licenseValid;
    }
}
