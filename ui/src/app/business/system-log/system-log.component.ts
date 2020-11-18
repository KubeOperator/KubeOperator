import {Component, OnInit} from '@angular/core';
import {SystemLogService} from './system-log.service';

@Component({
    selector: 'app-system-log',
    templateUrl: './system-log.component.html',
    styleUrls: ['./system-log.component.css']
})
export class SystemLogComponent implements OnInit {
    loading = false;
    total = 0;
    page = 1;
    size = 10;
    items = [];
    constructor(private service: SystemLogService) {}

    ngOnInit(): void {
        this.refresh()
    }
    refresh() {
        this.loading = true;
        this.service.list(this.page, this.size).subscribe(data => {
            this.items = data.items;
            this.total = data.total;
            this.loading = false;
        });
    }
}
