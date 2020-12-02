import {Component, OnInit} from '@angular/core';
import {BaseModelDirective} from "../../../shared/class/BaseModelDirective";
import {MultiClusterRepository} from "../multi-cluster-repository";
import {MultiClusterRepositoryService} from "../multi-cluster-repository.service";
import {Router} from "@angular/router";

@Component({
    selector: 'app-multi-cluster-repository-list',
    templateUrl: './multi-cluster-repository-list.component.html',
    styleUrls: ['./multi-cluster-repository-list.component.css']
})
export class MultiClusterRepositoryListComponent extends BaseModelDirective<MultiClusterRepository> implements OnInit {

    constructor(private multiClusterRepositoryService: MultiClusterRepositoryService, private router: Router) {
        super(multiClusterRepositoryService);
    }

    ngOnInit(): void {
        super.ngOnInit();
    }

    onDetail(name: string) {
        this.router.navigate(['multicluster', name]).then();
    }
}
