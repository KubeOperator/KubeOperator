import { ComponentFixture, TestBed } from '@angular/core/testing';

import { IpPoolCreateComponent } from './ip-pool-create.component';

describe('IpPoolCreateComponent', () => {
  let component: IpPoolCreateComponent;
  let fixture: ComponentFixture<IpPoolCreateComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ IpPoolCreateComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(IpPoolCreateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
