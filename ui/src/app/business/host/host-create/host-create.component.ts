import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {BaseModelDirective} from '../../../shared/class/BaseModelDirective';
import {Host, HostCreateRequest} from '../host';
import {HostService} from '../host.service';
import {NgForm} from '@angular/forms';
import {CredentialService} from '../../setting/credential/credential.service';
import {Credential} from '../../setting/credential/credential';
import {ModalAlertService} from '../../../shared/common-component/modal-alert/modal-alert.service';
import {AlertLevels} from '../../../layout/common-alert/alert';
import {CommonAlertService} from '../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {IpPattern} from '../../../constant/pattern';

@Component({
    selector: 'app-host-create',
    templateUrl: './host-create.component.html',
    styleUrls: ['./host-create.component.css']
})
export class HostCreateComponent extends BaseModelDirective<Host> implements OnInit {

    opened = false;
    isSubmitGoing = false;
    item: HostCreateRequest = new HostCreateRequest();
    ipPattern = IpPattern;
    credentials: Credential[] = [];
    @ViewChild('hostForm') hostForm: NgForm;
    @Output() created = new EventEmitter();

    constructor(private hostService: HostService, private credentialService: CredentialService,
                private modalAlertService: ModalAlertService,
                private commonAlertService: CommonAlertService, private translateService: TranslateService) {
        super(hostService);
    }

    ngOnInit(): void {
    }

    open() {
        this.credentialService.list().subscribe(data => {
            this.credentials = data.items;
        });
        this.opened = true;
        this.item = new HostCreateRequest();
        this.hostForm.resetForm({
            port: 22,
            credentialId: '',
        });
    }

    onCancel() {
        this.opened = false;
        this.item = new HostCreateRequest();
        this.hostForm.resetForm(this.item);
    }

    onSubmit() {
        this.isSubmitGoing = true;
        this.hostService.create(this.item).subscribe(data => {
            this.isSubmitGoing = false;
            this.created.emit();
            this.onCancel();
            this.commonAlertService.showAlert(this.translateService.instant('APP_ADD_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.isSubmitGoing = false;
            this.modalAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }
}
