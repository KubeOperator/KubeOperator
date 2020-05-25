import {Component, OnInit} from '@angular/core';
import {ActivatedRoute, Router} from '@angular/router';

@Component({
    selector: 'app-cluster-detail',
    templateUrl: './cluster-detail.component.html',
    styleUrls: ['./cluster-detail.component.css']
})
export class ClusterDetailComponent implements OnInit {

    constructor(private router: Router, private route: ActivatedRoute) {
    }

    ngOnInit(): void {
    }

    backToCluster() {
        this.router.navigate(['clusters']);
    }

}
