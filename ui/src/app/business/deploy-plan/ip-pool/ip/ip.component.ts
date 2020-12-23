import {Component, OnInit, ViewChild} from '@angular/core';
import {IpListComponent} from './ip-list/ip-list.component';
import {IpCreateComponent} from './ip-create/ip-create.component';
import {IpDeleteComponent} from './ip-delete/ip-delete.component';
import {ActivatedRoute, Router} from '@angular/router';

@Component({
    selector: 'app-ip',
    templateUrl: './ip.component.html',
    styleUrls: ['./ip.component.css']
})
export class IpComponent implements OnInit {

    ipPoolName: string;

    @ViewChild(IpListComponent, {static: true})
    list: IpListComponent;

    @ViewChild(IpCreateComponent, {static: true})
    create: IpCreateComponent;

    @ViewChild(IpDeleteComponent, {static: true})
    delete: IpDeleteComponent;

    constructor(private router: Router, private route: ActivatedRoute) {
        this.route.data.subscribe(data => {
            this.ipPoolName = data.ipPool.name;
        });
    }

    ngOnInit(): void {
    }

    backToIpPool() {
        this.router.navigate(['deploy/ip-pool']);
    }

    openCreate() {
        this.create.open();
    }

    openDelete(items) {
        this.delete.open(items);
    }

    refresh() {
        this.list.reset();
        this.list.refresh();
    }
}
