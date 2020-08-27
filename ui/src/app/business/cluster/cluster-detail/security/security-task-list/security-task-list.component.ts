import {Component, EventEmitter, Input, OnDestroy, OnInit, Output} from '@angular/core';
import {SecurityService} from "../security.service";
import {Cluster} from "../../../cluster";
import {CisTask} from "../security";

@Component({
    selector: 'app-security-task-list',
    templateUrl: './security-task-list.component.html',
    styleUrls: ['./security-task-list.component.css']
})
export class SecurityTaskListComponent implements OnInit, OnDestroy {

    constructor(private service: SecurityService) {
    }

    loading = false;
    items: CisTask[] = [];
    selected: CisTask[] = [];
    page = 1;
    size = 10;
    total = 0;
    @Output() detailEvent = new EventEmitter<CisTask>();
    @Output() createEvent = new EventEmitter();
    @Output() deleteEvent = new EventEmitter<CisTask[]>();
    @Input() currentCluster: Cluster;
    timer;


    ngOnInit(): void {
        this.refresh();
        this.polling();
    }

    ngOnDestroy() {
        clearInterval(this.timer);
    }

    refresh() {
        this.selected = [];
        this.service.page(this.currentCluster.name, this.page, this.size).subscribe(data => {
            this.items = data.items;
            this.total = data.total;
        });
    }

    onDetail(item: CisTask) {
        this.detailEvent.emit(item);
    }

    onCreate() {
        this.createEvent.emit();
    }

    onDelete() {
        this.deleteEvent.emit(this.selected);
    }

    polling() {
        this.timer = setInterval(() => {
            let flag = false;
            for (const item of this.items) {
                if (item.status === 'Running') {
                    flag = true;
                }
            }
            if (flag) {
                this.service.page(this.currentCluster.name, this.page, this.size).subscribe(data => {
                    this.items = data.items;
                    this.total = data.total;
                });
            }
        }, 5000);
    }

}
