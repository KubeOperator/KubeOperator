import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {Credential} from '../credential';
import {CredentialService} from '../credential.service';
import {AlertLevels} from '../../../../layout/common-alert/alert';
import {ModalAlertService} from '../../../../shared/common-component/modal-alert/modal-alert.service';
import {CommonAlertService} from '../../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';

@Component({
    selector: 'app-credential-delete',
    templateUrl: './credential-delete.component.html',
    styleUrls: ['./credential-delete.component.css']
})
export class CredentialDeleteComponent implements OnInit {

    opened = false;
    items: Credential[] = [];
    @Output() deleted = new EventEmitter();

    constructor(private service: CredentialService, private modalAlertService: ModalAlertService,
                private commonAlertService: CommonAlertService, private translateService: TranslateService) {
    }

    ngOnInit(): void {
    }


    open(items: Credential[]) {
        this.items = items;
        this.opened = true;
    }

    onCancel() {
        this.opened = false;
    }

    onSubmit() {
        this.service.batch('delete', this.items).subscribe(data => {
            this.deleted.emit();
            this.opened = false;
            this.commonAlertService.showAlert(this.translateService.instant('APP_DELETE_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.opened = false;
            this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }
}
