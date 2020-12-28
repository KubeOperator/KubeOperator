import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {ClusterService} from '../cluster.service';
import {Cluster} from '../cluster';
import {CommonAlertService} from '../../../layout/common-alert/common-alert.service';
import {AlertLevels} from '../../../layout/common-alert/alert';
import {TranslateService} from '@ngx-translate/core';

@Component({
    selector: 'app-cluster-delete',
    templateUrl: './cluster-delete.component.html',
    styleUrls: ['./cluster-delete.component.css']
})
export class ClusterDeleteComponent implements OnInit {

    opened = false;
    isSubmitGoing = false;
    items: Cluster[] = [];
    @Output() deleted = new EventEmitter();

    constructor(private service: ClusterService,private translateService: TranslateService, private commonAlert: CommonAlertService) {
    }


    ngOnInit(): void {
    }

    open(items: Cluster[]) {
        this.items = items;
        this.opened = true;
    }

    onCancel() {
        this.opened = false;
    }

    onSubmit() {
        if (this.isSubmitGoing) {
            return;
        }
        this.isSubmitGoing = true;
        this.opened = false;
        this.commonAlert.showAlert(this.translateService.instant('APP_DELETE_START'), AlertLevels.SUCCESS);
        this.service.batch('delete', this.items).subscribe(data => {
            this.deleted.emit();
            this.isSubmitGoing = false;
        }, error => {
            this.deleted.emit();
            this.isSubmitGoing = false;
            this.commonAlert.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }

}
