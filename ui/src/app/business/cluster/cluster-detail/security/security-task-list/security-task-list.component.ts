import {Component, OnInit} from '@angular/core';
import {SecurityTask} from "../security";

@Component({
    selector: 'app-security-task-list',
    templateUrl: './security-task-list.component.html',
    styleUrls: ['./security-task-list.component.css']
})
export class SecurityTaskListComponent implements OnInit {

    constructor() {
    }

    loading = false;
    selected: SecurityTask[] = [];
    items: SecurityTask[] = [];


    ngOnInit(): void {
    }

}
