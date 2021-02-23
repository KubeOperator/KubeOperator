import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {Manifest} from '../manifest';
import {ManifestService} from '../manifest.service';
import {CommonAlertService} from '../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {AlertLevels} from '../../../layout/common-alert/alert';

@Component({
    selector: 'app-manifest-list',
    templateUrl: './manifest-list.component.html',
    styleUrls: ['./manifest-list.component.css']
})
export class ManifestListComponent implements OnInit {
    constructor(private manifestService: ManifestService,
                private commonAlertService: CommonAlertService,
                private translateService: TranslateService) {
    }

    items: Manifest[] = [];
    loading = false;
    @Output() detailEvent = new EventEmitter<Manifest>();
    @Output() alertEvent = new EventEmitter();
    largeVersion = 'v1.20';


    ngOnInit(): void {
        this.refresh();
    }

    refresh() {
        this.manifestService.listByLargeVersion(this.largeVersion).subscribe(data => {
            this.items = data;
        });
    }

    changeVersion(largeVersion) {
        this.largeVersion = largeVersion;
        this.refresh();
    }

    onDetail(item: Manifest) {
        this.detailEvent.emit(item);
    }

    onAlert() {
        this.alertEvent.emit();
    }

    update(item: Manifest) {
        const updateItem = new Manifest();
        Object.assign(updateItem, item);
        updateItem.isActive = !item.isActive;
        if (updateItem.isActive) {
            this.onAlert();
        }
        this.manifestService.update(updateItem).subscribe(data => {
            this.commonAlertService.showAlert(this.translateService.instant('APP_UPDATE_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }
}
