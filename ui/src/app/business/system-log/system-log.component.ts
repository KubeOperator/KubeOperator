import {Component, OnInit, ViewChild} from '@angular/core';
import {SystemLogListComponent} from './system-log-list/system-log-list.component';

@Component({
    selector: 'app-system-log',
    templateUrl: './system-log.component.html',
    styleUrls: ['./system-log.component.css']
})
export class SystemLogComponent implements OnInit {
    @ViewChild(SystemLogListComponent)
    list: SystemLogListComponent;

    constructor(){}

    ngOnInit(): void {
    }
}
