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
        this.import.emit();
    }

    onSubmit() {
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
        let file = e.target.files[0]
        if (file.size > 10485760) {
            this.modalAlertService.showAlert(this.translateService.instant('APP_HOST_IMPORT_FILE_SIZE_ERROR'), AlertLevels.ERROR);
            return;
        }
        if (!this.endWith(file.name, 'xlsx') && !this.endWith(file.name, 'xls')) {
            this.modalAlertService.showAlert(this.translateService.instant('APP_HOST_IMPORT_FILE_ERROR'), AlertLevels.ERROR);
            return;
        }

        this.file = e.target.files[0];
    }

    endWith(str, suffix) {
        if(str == null || str == "" || suffix.length == 0 || suffix.length > str.length) {
            return false;
        }
        return str.substring(str.length - suffix.length) == suffix;
    }

    download() {
        window.open('/api/v1/hosts/template');
    }
}
