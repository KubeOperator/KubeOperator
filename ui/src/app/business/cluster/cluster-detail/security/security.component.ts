import {Component, OnInit, ViewChild} from '@angular/core';
import {ActivatedRoute} from "@angular/router";
import {Cluster} from "../../cluster";
import {SecurityTaskListComponent} from "./security-task-list/security-task-list.component";
import {SecurityTaskDetailComponent} from "./security-task-detail/security-task-detail.component";
import {CisTask} from "./security";
import {SecurityTaskCreateComponent} from "./security-task-create/security-task-create.component";
import {SecurityTaskDeleteComponent} from "./security-task-delete/security-task-delete.component";

@Component({
    selector: 'app-security',
    templateUrl: './security.component.html',
    styleUrls: ['./security.component.css']
})
export class SecurityComponent implements OnInit {

    constructor(private route: ActivatedRoute) {
    }

    currentCluster: Cluster;
    @ViewChild(SecurityTaskListComponent)
    list: SecurityTaskListComponent;
    @ViewChild(SecurityTaskDetailComponent)
    detail: SecurityTaskDetailComponent;
    @ViewChild(SecurityTaskCreateComponent)
    create: SecurityTaskCreateComponent;
    @ViewChild(SecurityTaskDeleteComponent)
    delete: SecurityTaskDeleteComponent;

    ngOnInit(): void {
        this.route.parent.data.subscribe(data => {
            this.currentCluster = data.cluster;
        });
    }

    openDetail(item: CisTask) {
        this.detail.open(item);
    }

    openCreate() {
        this.create.open();
    }

    refresh() {
        this.list.refresh();
    }

    openDelete(items: CisTask[]) {
        this.delete.open(items);
    }


}
