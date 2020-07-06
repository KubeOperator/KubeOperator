import {Component, OnInit} from '@angular/core';
import {Cluster} from "../../../cluster";
import {ActivatedRoute} from "@angular/router";

@Component({
    selector: 'app-chartmuseum',
    templateUrl: './chartmuseum.component.html',
    styleUrls: ['./chartmuseum.component.css']
})
export class ChartmuseumComponent implements OnInit {

    constructor(private route: ActivatedRoute) {
    }

    currentCluster: Cluster;

    ngOnInit(): void {
        this.route.parent.parent.data.subscribe(data => {
            this.currentCluster = data.cluster;
            console.log(this.currentCluster);
        });
    }

}
