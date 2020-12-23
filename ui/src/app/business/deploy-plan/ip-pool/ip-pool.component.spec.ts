import { ComponentFixture, TestBed } from '@angular/core/testing';

import { IpPoolComponent } from './ip-pool.component';

describe('IpPoolComponent', () => {
  let component: IpPoolComponent;
  let fixture: ComponentFixture<IpPoolComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ IpPoolComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(IpPoolComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
