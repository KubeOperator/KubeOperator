import {Component, OnInit} from '@angular/core';
import {BusinessLicenseService} from '../../shared/service/business-license.service';
import {LicenseService} from './license/license.service';

@Component({
    selector: 'app-setting',
    templateUrl: './setting.component.html',
    styleUrls: ['./setting.component.css']
})
export class SettingComponent implements OnInit {

    constructor(public businessLicenseService: BusinessLicenseService, private licenseService: LicenseService) {
    }

    hasLicense = false;

    ngOnInit(): void {
        this.licenseService.get().subscribe(data => {
            if (data.status === 'valid') {
                this.hasLicense = true;
            }
        });
        // this.hasLicense = this.businessLicenseService.licenseValid;
    }
}
