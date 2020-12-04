import {Component, OnInit} from '@angular/core';
import {ActivatedRoute, ActivatedRouteSnapshot, Router} from "@angular/router";
import {MultiClusterRepository} from "../multi-cluster-repository";

@Component({
    selector: 'app-multi-cluster-repository-detail',
    templateUrl: './multi-cluster-repository-detail.component.html',
    styleUrls: ['./multi-cluster-repository-detail.component.css']
})
export class MultiClusterRepositoryDetailComponent implements OnInit {

    constructor(private route: ActivatedRoute, private router: Router) {
    }

    currentRepository: MultiClusterRepository;

    ngOnInit(): void {
        this.route.data.subscribe(d => {
            this.currentRepository = d.repo;
        });
    }

    backToMultiCluster() {
        this.router.navigate(['multicluster']);
    }

}
