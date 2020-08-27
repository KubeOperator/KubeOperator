import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {CisTask} from "../security";
import {Cluster} from "../../../cluster";
import {CommonAlertService} from "../../../../../layout/common-alert/common-alert.service";
import {SecurityService} from "../security.service";
import {AlertLevels} from "../../../../../layout/common-alert/alert";

@Component({
    selector: 'app-security-task-delete',
    templateUrl: './security-task-delete.component.html',
    styleUrls: ['./security-task-delete.component.css']
})
export class SecurityTaskDeleteComponent implements OnInit {

    constructor(private alertService: CommonAlertService, private cisService: SecurityService) {
    }

    opened = false;
    items: CisTask[] = [];
    isSubmitGoing = false;
    @Input() currentCluster: Cluster;
    @Output() deleted = new EventEmitter();

    ngOnInit(): void {
    }

    open(items: CisTask[]) {
        this.items = items;
        this.opened = true;
    }

    onCancel() {
        this.opened = false;
    }

    onSubmit() {
        this.opened = false;
        const promises: Promise<{}>[] = [];
        this.items.forEach(item => {
            promises.push(this.cisService.delete(this.currentCluster.name, item.id).toPromise());
        });
        Promise.all(promises).then(() => {
        }, (error) => {
            this.alertService.showAlert(error, AlertLevels.ERROR);
        }).finally(() => {
            this.deleted.emit();
            this.opened = false;
        });
    }

}
