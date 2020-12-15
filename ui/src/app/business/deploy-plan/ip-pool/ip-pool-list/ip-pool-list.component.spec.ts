import { ComponentFixture, TestBed } from '@angular/core/testing';

import { IpPoolListComponent } from './ip-pool-list.component';

describe('IpPoolListComponent', () => {
  let component: IpPoolListComponent;
  let fixture: ComponentFixture<IpPoolListComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ IpPoolListComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(IpPoolListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
