import {Pipe, PipeTransform} from '@angular/core';

@Pipe({
    name: 'zoneStatus'
})
export class ZoneStatusPipe implements PipeTransform {

    transform(value: unknown, ...args: unknown[]): unknown {
        let result = '';
        if (value) {
            switch (value) {
                case 'READY':
                    result = '就绪';
                    break;
                case 'INITIALIZING':
                    result = '初始化中';
                    break;
                case 'UPLOADIMAGEERROR':
                    result = '上传镜像失败';
                    break;
            }
        }
        return result;
    }
}

