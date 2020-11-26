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
    @Output() import = new EventEmitter();

    constructor(private hostService: HostService, private modalAlertService: ModalAlertService,
                private commonAlertService: CommonAlertService, private translateService: TranslateService) {
    }

    ngOnInit(): void {
    }

    open() {
        this.opened = true;
    }

    onCancel() {
        this.opened = false;
    }

    onSubmit() {
        const formData = new FormData();
        formData.append('file', this.file);
        this.isSubmitGoing = true;
        this.hostService.upload(formData).subscribe(res => {
            // this.commonAlertService.showAlert(this.translateService.instant('APP_ADD_SUCCESS'), AlertLevels.SUCCESS);
            this.isSubmitGoing = false;
            console.log(res)
        }, error => {
            this.modalAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }

    upload(e) {
        this.file = e.target.files[0];
    }

    download() {
        window.open('/api/v1/hosts/template');
    }
}
