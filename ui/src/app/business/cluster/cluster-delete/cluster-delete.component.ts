import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {ClusterService} from '../cluster.service';
import {Cluster} from '../cluster';
import {CommonAlertService} from '../../../layout/common-alert/common-alert.service';
import {AlertLevels} from '../../../layout/common-alert/alert';

@Component({
    selector: 'app-cluster-delete',
    templateUrl: './cluster-delete.component.html',
    styleUrls: ['./cluster-delete.component.css']
})
export class ClusterDeleteComponent implements OnInit {

    opened = false;
    items: Cluster[] = [];
    @Output() deleted = new EventEmitter();

    constructor(private service: ClusterService, private commonAlert: CommonAlertService) {
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
        this.service.batch('delete', this.items).subscribe(data => {
            this.deleted.emit();
            this.opened = false;
        }, error => {
            this.deleted.emit();
            this.opened = false;
            this.commonAlert.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }

}
