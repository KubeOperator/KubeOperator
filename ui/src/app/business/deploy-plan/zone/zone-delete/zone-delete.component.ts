import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {BaseModelDirective} from '../../../../shared/class/BaseModelDirective';
import {Zone} from '../zone';
import {ModalAlertService} from '../../../../shared/common-component/modal-alert/modal-alert.service';
import {CommonAlertService} from '../../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {ZoneService} from '../zone.service';
import {AlertLevels} from '../../../../layout/common-alert/alert';

@Component({
    selector: 'app-zone-delete',
    templateUrl: './zone-delete.component.html',
    styleUrls: ['./zone-delete.component.css']
})
export class ZoneDeleteComponent extends BaseModelDirective<Zone> implements OnInit {

    opened = false;
    @Output() deleted = new EventEmitter();

    constructor(private zoneService: ZoneService, private modalAlertService: ModalAlertService,
                private commonAlertService: CommonAlertService, private translateService: TranslateService) {
        super(zoneService);
    }

    ngOnInit(): void {
    }


    open(items: Zone[]) {
        this.items = items;
        this.opened = true;
    }

    onCancel() {
        this.opened = false;
    }


    onSubmit() {
        this.zoneService.batch('delete', this.items).subscribe(data => {
            this.deleted.emit();
            this.opened = false;
            this.commonAlertService.showAlert(this.translateService.instant('APP_DELETE_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.opened = false;
            this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }
}
