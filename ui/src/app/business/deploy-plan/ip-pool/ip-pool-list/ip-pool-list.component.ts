import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {BaseModelDirective} from '../../../../shared/class/BaseModelDirective';
import {IpPool} from '../ip-pool';
import {IpPoolService} from '../ip-pool.service';

@Component({
    selector: 'app-ip-pool-list',
    templateUrl: './ip-pool-list.component.html',
    styleUrls: ['./ip-pool-list.component.css']
})
export class IpPoolListComponent extends BaseModelDirective<IpPool> implements OnInit {

    constructor(private ipPoolService: IpPoolService) {
        super(ipPoolService);
    }

    ngOnInit(): void {
        super.ngOnInit();
    }

}
