import {Component, Input, OnInit} from '@angular/core';
import {ChartmuseumService} from "../chartmuseum.service";
import {Chart} from "../chart";
import {Cluster} from "../../../../cluster";

@Component({
    selector: 'app-chart-list',
    templateUrl: './chart-list.component.html',
    styleUrls: ['./chart-list.component.css']
})
export class ChartListComponent implements OnInit {

    loading = false;
    items: Chart[] = [];
    @Input() currentCluster: Cluster;

    constructor(private chartService: ChartmuseumService) {
    }

    ngOnInit(): void {
        this.refresh();
    }

    refresh() {
        this.loading = true;
        this.chartService.list(this.currentCluster.name).subscribe(data => {
            this.items = [];
            for (const key in data) {
                if (key) {
                    this.items = this.items.concat(data[key]);
                }
            }
            this.loading = false;
        });
    }
}
