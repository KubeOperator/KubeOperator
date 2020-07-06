import {Component, OnInit} from '@angular/core';
import {ActivatedRoute} from "@angular/router";
import {Cluster} from "../../../cluster";

@Component({
    selector: 'app-registry',
    templateUrl: './registry.component.html',
    styleUrls: ['./registry.component.css']
})
export class RegistryComponent implements OnInit {

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
