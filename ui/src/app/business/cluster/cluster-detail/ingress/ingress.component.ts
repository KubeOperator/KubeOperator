import {Component, OnInit} from '@angular/core';
import {Cluster} from '../../cluster';
import {ActivatedRoute} from '@angular/router';

@Component({
    selector: 'app-ingress',
    templateUrl: './ingress.component.html',
    styleUrls: ['./ingress.component.css']
})
export class IngressComponent implements OnInit {

    constructor(private route: ActivatedRoute) {
    }

    currentCluster: Cluster;

    ngOnInit(): void {
        this.route.parent.data.subscribe(data => {
            this.currentCluster = data.cluster.item;
        });
    }

}
