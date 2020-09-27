import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {BaseModelDirective} from '../../../../shared/class/BaseModelDirective';
import {Region} from '../region';
import {RegionService} from '../region.service';
import {ModalAlertService} from '../../../../shared/common-component/modal-alert/modal-alert.service';
import {CommonAlertService} from '../../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {AlertLevels} from '../../../../layout/common-alert/alert';

@Component({
    selector: 'app-region-delete',
    templateUrl: './region-delete.component.html',
    styleUrls: ['./region-delete.component.css']
})
export class RegionDeleteComponent extends BaseModelDirective<Region> implements OnInit {

    opened = false;
    items: Region[] = [];
    @Output() deleted = new EventEmitter();

    constructor(private regionService: RegionService, private modalAlertService: ModalAlertService,
                private commonAlertService: CommonAlertService, private translateService: TranslateService) {
        super(regionService);
    }

    ngOnInit(): void {
    }

    open(items: Region[]) {
        this.items = items;
        this.opened = true;
    }

    onCancel() {
        this.opened = false;
    }

    onSubmit() {
        this.regionService.batch('delete', this.items).subscribe(data => {
            this.deleted.emit();
            this.opened = false;
            this.commonAlertService.showAlert(this.translateService.instant('APP_DELETE_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
            this.opened = false;
        });
    }
}
