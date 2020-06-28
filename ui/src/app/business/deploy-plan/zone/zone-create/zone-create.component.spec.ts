import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ZoneCreateComponent } from './zone-create.component';

describe('ZoneCreateComponent', () => {
  let component: ZoneCreateComponent;
  let fixture: ComponentFixture<ZoneCreateComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ZoneCreateComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ZoneCreateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
