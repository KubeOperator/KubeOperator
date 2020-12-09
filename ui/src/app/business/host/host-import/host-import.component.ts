import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {HostService} from '../host.service';
import {ModalAlertService} from '../../../shared/common-component/modal-alert/modal-alert.service';
import {CommonAlertService} from '../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {AlertLevels} from '../../../layout/common-alert/alert';

@Component({
    selector: 'app-host-import',
    templateUrl: './host-import.component.html',
    styleUrls: ['./host-import.component.css']
})
export class HostImportComponent implements OnInit {


    opened = false;
    isSubmitGoing = false;
    file;
    errMsg = '';
    @Output() import = new EventEmitter();

    constructor(private hostService: HostService, private modalAlertService: ModalAlertService,
                private commonAlertService: CommonAlertService, private translateService: TranslateService) {
    }

    ngOnInit(): void {
    }

    open() {
        this.opened = true;
        this.isSubmitGoing = false;
        this.errMsg = '';
    }

    onCancel() {
        this.opened = false;
        this.errMsg = '';
    }

    onSubmit() {
        const startIndex = this.file.name.lastIndexOf('.');
        if (startIndex !== -1) {
            const fileType = this.file.name.substring(startIndex + 1, this.file.name.length).toLowerCase();
            if (fileType !== 'xlsx') {
                this.modalAlertService.showAlert(this.translateService.instant('APP_HOST_IMPORT_FILE_ERROR'), AlertLevels.ERROR);
                return;
            }
        } else {
            this.modalAlertService.showAlert(this.translateService.instant('APP_HOST_IMPORT_FILE_ERROR'), AlertLevels.ERROR);
            return;
        }
        const formData = new FormData();
        formData.append('file', this.file);
        this.isSubmitGoing = true;
        this.hostService.upload(formData).subscribe(res => {
            this.isSubmitGoing = false;
            this.commonAlertService.showAlert(this.translateService.instant('APP_ADD_SUCCESS'), AlertLevels.SUCCESS);
            this.opened = false;
            this.import.emit();
        }, error => {
            this.isSubmitGoing = false;
            this.errMsg = error.error.msg;
        });
    }

    upload(e) {
        this.file = e.target.files[0];
    }

    download() {
        window.open('/api/v1/hosts/template');
    }
}
