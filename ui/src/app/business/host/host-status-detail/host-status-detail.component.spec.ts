import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { HostStatusDetailComponent } from './host-status-detail.component';

describe('HostStatusDetailComponent', () => {
  let component: HostStatusDetailComponent;
  let fixture: ComponentFixture<HostStatusDetailComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ HostStatusDetailComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(HostStatusDetailComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
