import {Component, OnInit, ViewChild} from '@angular/core';
import {LicenseImportComponent} from './license-import/license-import.component';
import {LicenseService} from './license.service';
import {License} from './license';
import {CommonAlertService} from '../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {AlertLevels} from '../../../layout/common-alert/alert';

@Component({
    selector: 'app-license',
    templateUrl: './license.component.html',
    styleUrls: ['./license.component.css']
})
export class LicenseComponent implements OnInit {
    licenseInfo: License = new License();
    licenseStatus = '';
    constructor(private licenseService: LicenseService,
                private commonAlertService: CommonAlertService,
                private translateService: TranslateService) {
    }

    licDate = {
        isDuringDate(endDateStr) {
            const c = new Date();
            const curDate = new Date(c.getTime() + 168 * 60 * 60 * 1000);
            const endDate = new Date(endDateStr);
            if (curDate >= endDate) {
                return c < endDate;
            }
            return false;
        }
    };

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
            this.licenseInfo.license = data.license;
            console.log(data);
            if ( data.status !== '') {
                this.licenseStatus = data.status;
            }
            if (this.licDate.isDuringDate(data.license.expired)) {
                this.commonAlertService.showAlert(this.translateService.instant('APP_LICENSE_EXPIRED_MSG'), AlertLevels.ERROR);
            }
        });
    }

}
