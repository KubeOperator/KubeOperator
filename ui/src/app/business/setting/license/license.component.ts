import {Component, OnInit, ViewChild} from '@angular/core';
import {LicenseImportComponent} from "./license-import/license-import.component";
import {LicenseService} from "./license.service";
import {License} from "./license";

@Component({
    selector: 'app-license',
    templateUrl: './license.component.html',
    styleUrls: ['./license.component.css']
})
export class LicenseComponent implements OnInit {
    license: License = new License();
    licenseStatus = 'invalid'
    constructor(private licenseService: LicenseService) {
    }

    @ViewChild(LicenseImportComponent, {static: true})
    import: LicenseImportComponent;

    ngOnInit(): void {
        this.refresh();
    }

    onImport() {
        this.import.open();
    }

    refresh() {
        this.licenseService.get().subscribe(data => {
            this.license = data;
            this.licenseStatus = this.license.status
        });
    }

}
