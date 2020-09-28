import {Component, OnInit} from '@angular/core';
import {LicenseService} from '../../business/setting/license/license.service';

@Component({
    selector: 'app-navigation',
    templateUrl: './navigation.component.html',
    styleUrls: ['./navigation.component.css']
})
export class NavigationComponent implements OnInit {

    constructor(private licenseService: LicenseService) {
    }
    hasLicense = false;

    ngOnInit(): void {
        this.licenseService.get().subscribe(data => {
            this.hasLicense = true;
        });
    }

}
