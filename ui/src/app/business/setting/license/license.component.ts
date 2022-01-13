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
    license: License = new License();
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
        this.licenseService.gethw().subscribe(data => {
            this.license = data;
            if ( this.license.status !== '') {
                this.licenseStatus = this.license.status;
            }
            if (this.licDate.isDuringDate(data.expired)) {
                this.commonAlertService.showAlert(this.translateService.instant('APP_LICENSE_EXPIRED_MSG'), AlertLevels.ERROR);
            }
        });
    }

}
