import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {BaseModelDirective} from '../../../../../shared/class/BaseModelDirective';
import {Ip} from '../ip';
import {IpService} from '../ip.service';
import {AlertLevels} from '../../../../../layout/common-alert/alert';
import {ModalAlertService} from '../../../../../shared/common-component/modal-alert/modal-alert.service';
import {CommonAlertService} from '../../../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {ActivatedRoute, Router} from '@angular/router';

@Component({
    selector: 'app-ip-delete',
    templateUrl: './ip-delete.component.html',
    styleUrls: ['./ip-delete.component.css']
})
export class IpDeleteComponent extends BaseModelDirective<Ip> implements OnInit {

    opened = false;
    isSubmitGoing = false;
    items: Ip[] = [];
    ipPoolName = '';
    @Output() deleted = new EventEmitter();

    constructor(private ipService: IpService,
                private modalAlertService: ModalAlertService,
                private commonAlertService: CommonAlertService,
                private translateService: TranslateService,
                private router: Router,
                private route: ActivatedRoute) {
        super(ipService);
    }

    ngOnInit(): void {
        this.route.data.subscribe(data => {
            this.ipPoolName = data.ipPool.name;
        });
    }

    open(items) {
        this.items = items;
        this.opened = true;
    }

    onCancel() {
        this.isSubmitGoing = false;
        this.opened = false;
    }

    onSubmit() {
        this.isSubmitGoing = true;
        this.ipService.batch('delete', this.items, this.ipPoolName).subscribe(res => {
            this.isSubmitGoing = false;
            this.opened = false;
            this.deleted.emit();
            this.commonAlertService.showAlert(this.translateService.instant('APP_DELETE_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.isSubmitGoing = false;
            this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
            this.opened = false;
        });
    }
}
