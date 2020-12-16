import {Component, OnInit} from '@angular/core';
import {BaseModelDirective} from '../../../../../shared/class/BaseModelDirective';
import {Ip} from '../ip';
import {IpService} from '../ip.service';
import {ActivatedRoute, Router} from '@angular/router';
import {IpPool} from '../../ip-pool';

@Component({
    selector: 'app-ip-list',
    templateUrl: './ip-list.component.html',
    styleUrls: ['./ip-list.component.css']
})
export class IpListComponent extends BaseModelDirective<Ip> implements OnInit {

    ipPoolName: string;
    ipPool: IpPool;

    constructor(private ipService: IpService,
                private router: Router,
                private route: ActivatedRoute) {
        super(ipService);
    }

    ngOnInit(): void {
        this.route.data.subscribe(data => {
            this.ipPool = data.ipPool;
            this.ipPoolName = data.ipPool.name;
            this.refresh();
        });
    }

    refresh() {
        this.loading = true;
        this.ipService.page(this.page, this.size, this.ipPoolName).subscribe(data => {
            this.items = data.items;
            this.total = data.total;
            this.loading = false;
        });
    }
}
