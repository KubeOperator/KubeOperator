import {Component, OnInit, ViewChild} from '@angular/core';
import {PlanListComponent} from './plan-list/plan-list.component';
import {PlanCreateComponent} from './plan-create/plan-create.component';
import {PlanDeleteComponent} from './plan-delete/plan-delete.component';

@Component({
    selector: 'app-plan',
    templateUrl: './plan.component.html',
    styleUrls: ['./plan.component.css']
})
export class PlanComponent implements OnInit {

    @ViewChild(PlanListComponent, {static: true})
    list: PlanListComponent;

    @ViewChild(PlanCreateComponent, {static: true})
    create: PlanCreateComponent;

    @ViewChild(PlanDeleteComponent, {static: true})
    delete: PlanDeleteComponent;

    constructor() {
    }

    ngOnInit(): void {
    }

    refresh() {
        this.list.reset();
        this.list.refresh();
    }

    openCreate() {
        this.create.open();
    }

    openDelete(items) {
        this.delete.open(items);
    }
}
