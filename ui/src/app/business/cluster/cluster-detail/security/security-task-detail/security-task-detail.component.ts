import {Component, OnInit} from '@angular/core';
import {CisTask} from "../security";

@Component({
    selector: 'app-security-task-detail',
    templateUrl: './security-task-detail.component.html',
    styleUrls: ['./security-task-detail.component.css']
})
export class SecurityTaskDetailComponent implements OnInit {

    constructor() {
    }

    item: CisTask = new CisTask();
    opened = false;

    ngOnInit(): void {
    }


    open(item: CisTask) {
        this.opened = true;
        this.item = item;
    }

}
